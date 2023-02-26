#!/usr/bin/env ruby
# frozen_string_literal: true

RED     = "\033[31m"
GREEN   = "\033[32m"
YELLOW  = "\033[33m"
BLUE    = "\033[34m"
MAGENTA = "\033[35m"
CYAN    = "\033[36m"
RESET   = "\033[0m"

vars = []

while (line = ARGF.gets)
  pos = 0
  loop do
    m = /\$\{(\w+)(:[-=](.*))?\}/.match(line, pos) || /\$(\w+)/.match(line, pos)

    break unless m

    name = m[1]
    default_val = m[3]
    vars << [name, ARGF.lineno, default_val]
    pos = m.offset(0).last
  end
end

UNDEFINED = '<undefined>'

name_max_len = vars.map { |(name, _, _)| name.length }.max
val_max_len = vars.map { |(_, _, v)| v ? v.length : 0 }.max
val_max_len = [val_max_len, UNDEFINED.length].max

puts "#{'name'.center(name_max_len + 1, '-')} #{'default'.center(val_max_len, '-')} #{'lineno'.center(10, '-')}"

vars.sort!.each do |(name, lineno, default_val)|
  if default_val
    val = default_val
    width = val_max_len
  else
    val = "#{MAGENTA}#{UNDEFINED}"
    width = val_max_len + MAGENTA.length
  end
  printf "#{CYAN}$%-#{name_max_len}s#{RESET} %-#{width}s#{RESET} (line#{YELLOW}%3d#{RESET})\n", name, val, lineno
end
