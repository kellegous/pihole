require 'fileutils'
require './build'

ENV['PATH'] = "#{Dir.pwd}/bin:#{ENV['PATH']}"

GOOS = ENV['GOOS'] || `go env GOOS`.strip
GOARCH = ENV['GOARCH'] || `go env GOARCH`.strip
DESC = "#{GOOS}-#{GOARCH}"
puts DESC

set_gopath(['.'])

GODEPS = go_get('src', [
	'github.com/BurntSushi/toml',
	'golang.org/x/crypto/ssh',
	'golang.org/x/crypto/ssh/agent',
	'github.com/golang/protobuf/...',
	'google.golang.org/grpc',
	'github.com/golang/glog'
])

PROTOS = protoc('src/pihole')

DEPS = GODEPS + PROTOS

TARGS = [
	"bin/#{DESC}/pihole",
	"bin/#{DESC}/piholed"
]

task :subl do
	sh 'subl', '.'
end

file "bin/#{DESC}/pihole" => DEPS + FileList['src/pihole/**/*'] do |t|
	sh 'go', 'build', '-o', t.name, 'pihole/client'
end

file "bin/#{DESC}/piholed" => DEPS + FileList['src/pihole/**/*'] do |t|
	sh 'go', 'build', '-o', t.name, 'pihole/server'
end

task :default => TARGS

task :test do
	sh 'go', 'test', 'pihole/...'
end

task :clean do
	TARGS.each do |f|
		FileUtils.rm(f)
	end
end
