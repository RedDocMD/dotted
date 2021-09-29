#!/usr/bin/env ruby

image_name = "redocmd/dotted"
homedir = Dir.home()
pwd = Dir.pwd()

volume_codedir = "/code"
volume_godir = "/godir"
default_godir = File.join(homedir, "go")

build_command = "docker build -t #{image_name} ."

if ARGV.size != 1 && ARGV.size != 3
    puts "Expected: [--godir GODIR] [build | run]"
    exit(1)
end

command = ARGV[-1]

if command != "build" && command != "run"
    puts "Expected command to be \"build\" or \"run\""
    exit(1)
end

def color(str, color)
    if color == 'red'
        return "\033[31m#{str}\033[0m"
    elsif color == 'green'
        return "\033[32m#{str}\033[0m"
    end
end

def color_bold(str, color)
    if color == 'red'
        return "\033[31;1m#{str}\033[0m"
    elsif color == 'green'
        return "\033[32;1m#{str}\033[0m"
    end
end

if command == "build"
    puts "Command: #{color(build_command, "red")}"
    puts color_bold("Building ...", "green")
    system build_command
else
    if ARGV.size == 3
        if ARGV[1] != "--godir"
            puts "Expected arg to be --godir"
            exit(1)
        end
        godir = ARGV[1]
    else
        godir = default_godir
    end
    run_command = "docker run --rm -it -v #{pwd}:#{volume_codedir} -v #{godir}:#{volume_godir} #{image_name}"
    puts "Command: #{color(run_command, "red")}"
    puts color_bold("Running ...", "green")
    system run_command
end