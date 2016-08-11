
def set_gopath(paths)
  ENV['GOPATH'] = paths.map { |p|
    "#{Dir.pwd}/#{p}"
  }.join(':')
end

def go_get(dst, deps)
	deps.map do |pkg|
		path = pkg.gsub(/\/\.\.\.$/, '')
		dest = File.join(dst, path)
		file dest do
			sh 'go', 'get', pkg
		end
		dest
	end
end

class Path
  def <<(p)
    path = ENV['PATH']
    p = File.join(Dir.pwd, p)
    ENV['PATH'] = "#{p}:#{path}"
  end
end

def path
  Path.new
end

# generates build rules for protobufs. Rules that target dst are generated
# by scanning src.
def protoc(src)
  FileList["#{src}/**/*.proto"].map do |src_path|
    dst_path = src_path.sub(/\.proto/, '.pb.go')
    file dst_path => [src_path] do
      sh 'protoc', "-I#{src}", src_path, "--go_out=plugins=grpc:#{src}"
    end

    dst_path
  end
end
