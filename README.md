# go-gnome-session-inhibit

It's a library that allows to acquire (and release) session inhibits from Gnome Session Manager
(using DBus for communication) in order to prevent operating system from going idle (screen off)
or suspend/sleep (also logout, switching-user and automount because that's all 5 operations
supported by Gnome Session Manager).

Library is really simple - please just see included example. I've written it for use in my other
tools and to have code in one place, but may be useful for other folks ;).
