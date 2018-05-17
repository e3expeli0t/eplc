def gen_for_one(string):
	pname, prod = string.split("::=")
	productions = prod.split("|")

	for x in productions:
		 print("{0}::={1}".format(pname, x))

def gen(fname):
	with open(fname, "r") as reader:
		to_ret = []
		for line in reader:
			x = gen_for_one(line.strip()) #I know it's super ineffiecient but i'm lazy...
			if x != None:
				print (x)

gen("grammer.txt")				