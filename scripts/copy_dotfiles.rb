require 'fileutils'
require 'pathname'

loc_file = "/root/data/loc.dat"
data_dir = "/root/data"

File.foreach(loc_file) do |line|
    parts = line.split(" ")
    name = parts[0]
    src = File.join(data_dir, name)
    dest_rel = parts[1]
    dest = File.join("/root", dest_rel)
    parent = Pathname.new(dest).parent
    FileUtils.mkdir_p(parent)
    FileUtils.cp(src, dest)
end