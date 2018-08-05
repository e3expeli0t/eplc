
#ifndef EDBG_DEBUGGER_HPP
#define EDBG_DEBUGGER_HPP

#include <utility>
#include <string>
#include <linux/types.h>


class debugger {
 public:
     debugger (std::string prog_name, pid_t pid)
         : m_prog_name{std::move(prog_name)}, m_pid{pid} {}
     void run();
 private:
     void handle_command(const std::string& line);
     void continue_execution();        
     
     std::string m_prog_name;
     pid_t m_pid;
 };

class breakpoint {
public:
    breakpoint(pid_t pid, std::intptr_t addr)
        : m_pid{pid}, m_addr{addr}, m_enabled{false}, m_saved_data{}
    {}

    void enable();
    void disable();

    auto is_enabled() const -> bool { return m_enabled; }
    auto get_address() const -> std::intptr_t { return m_addr; }

private:
    pid_t m_pid;
    std::intptr_t m_addr;
    bool m_enabled;
    uint8_t m_saved_data; //data which used to be at the breakpoint address
};

#endif