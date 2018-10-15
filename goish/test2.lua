a, b = print(1,"hi")
print(1 and 0 or 3)
print(5+6)

inc = function(a)
	return a+1
end

print(inc(42))

print(a == 1 and 7 or 9)

t = {
	a = 1,
	"__*&^" = 2,
	c = {
		d = function()
			return "hello world"
		end
	}
}

print(t.c.d())
t.a = 42
print(t)

print("\n", t:keys())

print(2*2+3)

t = {
	g = function()
		return "hello world"
	end,
}

t.m = function(self)
	return self:g()
end

print(t:m())

s = "goodbye world"

print(s:len())
print(s:type())

print([1,2,7])

print("\n")

a = {}
print(a, "\n")
a:set("1", 1)
print(a, "\n")

b = {}
a:set(b, 2)
print(a, "\n")
b:set("2", 2)
print(a, "\n")

l = [1, 2, 3]
print(l, "\n")
l:push(4)
print(l, "\n")
print(l:pop())
print(l, "\n")

print("a" + "b", "\n")

len = "hi"
print("yo", l:len(), "\n")

print("a,b,c":split(","))
print("a,b,c":split(","):join(":"))

c = chan(10)
c:write(1)
c:write(2)
c:write(3)

print(c:read())
print(c:read())
print(c:read())

hi = function(g)
	print("hi\n")
end

g = group()
g:run(hi)
g:run(hi)
g:run(hi)
g:wait()
print("done\n")

print(b:get("hi"))

print(true and "yes" or "no")

true = 0 < 1
false = 1 < 0
nil = []:pop()

for 10 do(i)
	print(i, "\n")
end

list = []:lib()

print(list, "\n")

list.iterate = function(self)
	i = 0
	return do
		l = self:len()
		if i < l then
			i = i+1
			print("len ", l, " ", i, " ", self:get(i-1), "\n")
			return i, self:get(i-1)
		end
	end
end

print(list, "\n")
fn = [1,2,3]:iterate()
print(fn(), "\n")
print(fn(), "\n")
print(fn(), "\n")

--for [1,2,3] do(i, v)
--	print(i, ":", v, "\n")
--end
