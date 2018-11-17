function arrayish(a) {
  try {
    return a.length > 0 && '0' in Object(a);
  }
  catch(e) {
    return false;
  }
}

function foreach(t, f) {
	if (typeof t === 'null' || typeof t === 'undefined') {
		return;
	}
	if (arrayish(t)) {
		for (var i = 0; i < t.length; i++) {
			f(i, t[i]);
		}
		return;
	}
	if (typeof t === 'object') {
		for (var k in t) {
			if (t.hasOwnProperty(k)) {
				f(k, t[k]);
			}
		};
		return;
	}
	f(0, t);
}

function post(url, vars, fn) {
	var fd = new FormData();
	foreach(vars, function(k, v) {
		fd.append(k, v);
	});
	window.fetch(url, {
		method: 'POST',
		body: fd,
	})
	.catch(error => console.error(error))
	.then(response => response.json())
	.then(fn);
}

function select(parent, sel) {
	var items = [];
	if (sel == undefined) {
		items.push(parent);
	} else {
		var list = parent.querySelectorAll(':scope '+sel);
		for (var i = 0; i < list.length; i++) {
			items.push(list[i]);
		};
	}
	var r = {};
	r.all = function() {
		return items;
	}
	r.one = function() {
		return items[0];
	}
	r.css = function(style) {
		foreach(items, (i, item) => {
			foreach(style, (k, v) => {
				if (v == null) {
					delete item.style[k];
				} else {
					item.style[k] = v;
				}
			});
		});
		return r;
	}
	r.hide = function() {
		r.css({ display: 'none' })
		return r;
	}
	r.show = function() {
		r.css({ display: null })
		return r;
	}
	r.remove = function() {
		foreach(items, (i, item) => {
			item.remove();
		});
		return r;
	}
	r.append = function(child) {
		foreach(items, (i, item) => {
			item.appendChild(child);
		});
		return r;
	}
	r.html = function(html) {
		foreach(items, (i, item) => {
			item.innerHTML = html;
		});
		return r;
	}
	r.text = function(text) {
		foreach(items, (i, item) => {
			item.innerText = text;
		});
		return r;
	}
	r.attr = function(key, val) {
		foreach(items, (i, item) => {
			item[key] = val;
		});
		return r;
	}
	r.tooltip = function(val) {
		foreach(items, (i, item) => {
			item.title = val;
		});
		return r;
	}
	r.click = function(fn) {
		foreach(items, (i, item) => {
			item.addEventListener('click', fn);
		});
		return r;
	}
	r.clone = function() {
		foreach(items, (i, item) => {
			items[i] = item.cloneNode(true);
		});
		return r;
	}
	r.hotkey = function(combo, cb) {
		var keys = {};
		foreach(combo.trim().split(/\s/), (i, key) => {
			if (key.toLowerCase() == 'space')
				key = ' ';
			keys[key] = false;
		});
		var handler = (e, fn) => {
			if (keys.hasOwnProperty(e.key)) {
				fn(e.key);
				e.preventDefault();
				e.stopPropagation();
			}
		}
		var groups = [];
		foreach(items, (i, item) => {
			var up = function(e) {
				handler(e, (key) => {
					var run = true;
					foreach(keys, (k, state) => {
						run = state;
					});
					if (run) {
						cb();
					}
					keys[key] = false;
				});
			};
			var down = function(e) {
				handler(e, (key) => {
					keys[key] = true;
				});
			};
			groups.push([item, up, down]);
			item.addEventListener('keyup', up);
			item.addEventListener('keydown', down);
		});
		return function() {
			foreach(groups, (i, group) => {
				var item = group[0];
				item.removeEventListener('keyup', group[1]);
				item.removeEventListener('keydown', group[2]);
			});
		}
	}

	return r;
}

function create(tag) {
	return document.createElement(tag);
}
