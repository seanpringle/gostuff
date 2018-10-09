a, b = print(1,"hi")
print(1 and 0 or 3)
print(5+6)

inc = function(a) do return a+1 end
print(inc(42))

print(a == 1 and 7 or 9)

t = { a = 1, "__*&^" = 2, c = { d = function() do return "hello world" end}}

print(t.c.d())
t.a = 42
print(t)

print("\n", t:keys())

print(2*2+3)

t = {
	g = function() do
		return "hello world"
	end,
}

t.m = function(self) do
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
print(a:getmeta(), "\n")

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

hi = function(g) do
	print("hi\n")
end

g = group()
g:run(hi)
g:run(hi)
g:run(hi)
g:wait()
print("done\n")

print(b:get("hi"))

true = 0 < 1
false = 1 < 0
nil = []:pop()

iter = function(list) do
	t = { pos = 0 }
	return function() do
		v = list:get(t.pos)
		t.pos = t.pos + 1
		return v
	end
end

for it = iter([1,2,3]); i = it() do
	print(i, "\n")
end

print("\n")

print(true and "yes" or "no")

for i = 0; i < 10; i = i + 1 do
	print(i, "\n")
end

print(
	function() do return [1,2,3] end(),
	function() do return [4,5,6] end(),
)
