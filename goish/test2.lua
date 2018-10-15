print(a, b = print(1,"hi"))
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

print("", t:keys())

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
print(type(s))

print([1,2,7])

a = {}
print(a)
a:set("1", 1)
print(a)

b = {}
a:set(b, 2)
print(a)
b:set("2", 2)
print(a)

l = [1, 2, 3]
print(l)
l:push(4)
print(l)
print(l:pop())
print(l)

print("a" + "b")

len = "hi"
print("yo", l:len())

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
	print("hi")
end

g = group()
g:run(hi)
g:run(hi)
g:run(hi)
g:wait()
print("done")

print(b:get("hi"))

print(true and "yes" or "no")

true = 0 < 1
false = 1 < 0
nil = []:pop()

function(proto)
	proto.iterate = function(self)
		i = 0
		return do
			if i < self then
				return i++
			end
		end
	end
end(getprototype(0))

for 10 do(i)
	print(i)
end

function(list)
	list.iterate = function(self)
		i = 0
		return do
			if i < self:len() then
				return i++, self:get(i-1)
			end
		end
	end
end(getprototype([]))

for [1,2,3] function(i, v)
	print(i, ":", v)
end

function(map)
	map.iterate = function(self)
		i = 0
		keys = self:keys()
		return do
			if i < keys:len() then
				key = keys:get(i++)
				return key, self:get(key)
			end
		end
	end
end(getprototype({}))

for { tom = 1, dick = 2, harry = 43 } function(k, v)
	print(k, "=>", v)
end

a = 1
print(a++)
print(a++)
print(a++)

for 10 do(i)
	if i == 5 then
		break
	end
	print(i)
end

