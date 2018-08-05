#include <iostream>
#include <vector>
#include <sstream>
#include <cstring>
#include <sys/wait.h>
#include <sys/ptrace.h>
#include <unistd.h>

#include "luncher.h"
#include "linenoise.h"


std::vector<std::string> split(const std::string &s, char delimiter) {
    std::vector<std::string> out{};
    std::stringstream ss {s};
    std::string item;

    while (std::getline(ss,item,delimiter)) {
        out.push_back(item);
    }

    return out;
}

bool is_prefix(const std::string& s, const std::string& of) {
    if (s.size() > of.size()) return false;
    return std::equal(s.begin(), s.end(), of.begin());
}


void debugger::run() {
    int wait_status;
    auto options = 0;
    waitpid(m_pid, &wait_status, options);

    char* line = nullptr;
    while((line = linenoise("epldbg-> ")) != nullptr) {
        handle_command(line);
        linenoiseHistoryAdd(line);
        linenoiseFree(line);
    }
}

void debugger::handle_command(const std::string& line) {
    auto args = split(line,' ');
    auto command = args[0];

    if (is_prefix(command, "continue")) {
        continue_execution();
    } else {
        std::cerr << "Unknown command '" << command << "' \n";
    }
}

void debugger::continue_execution() {
    ptrace(PTRACE_CONT, m_pid, nullptr, nullptr);

    int wait_status;
    auto options = 0;
    waitpid(m_pid, &wait_status, options);
}

void breakpoint::enable() {
    auto data = ptrace(PTRACE_PEEKDATA, m_pid, m_addr, nullptr);
    m_saved_data = static_cast<uint8_t>(data & 0xff); //save bottom byte
    uint64_t int3 = 0xcc;
    uint64_t data_with_int3 = ((data & ~0xff) | int3); //set bottom byte to 0xcc
    ptrace(PTRACE_POKEDATA, m_pid, m_addr, data_with_int3);

    m_enabled = true;
}

int main(int argc, char** argv) {
    if (argc < 2) 
    {
        std::cerr << "Usage: eplsbg <program> \n";
        return -1;
    }

    auto prog = argv[1];

    auto pid = fork();
    if (pid == 0)  
    {
        //we're in the child process
        //execute debugee
           if (ptrace(PTRACE_TRACEME, 0, 0, 0) < 0) {
                std::cerr << "[!] Error in ptrace\n";
                return 1;
           }
           execl((char*)prog, ((char*)prog), nullptr);
    }
    else if (pid >= 1)  {
        //we're in the parent process
        //execute debugger
        debugger edbg{prog, pid};
        edbg.run();
    }
}
