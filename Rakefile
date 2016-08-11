require 'fileutils'
require './build'

ENV['PATH'] = "#{Dir.pwd}/bin:#{ENV['PATH']}"

set_gopath(['.'])

GODEPS = go_get('src', [
	'github.com/BurntSushi/toml',
	'golang.org/x/crypto/ssh',
	'golang.org/x/crypto/ssh/agent',
	'github.com/golang/protobuf/...',
	'google.golang.org/grpc',
])

PROTOS = protoc('src/pihole')

DEPS = GODEPS + PROTOS

task :atom do
	sh 'atom', '.'
end

task :subl do
	sh 'subl', '.'
end

file 'bin/pihole' => DEPS + FileList['src/pihole/**/*'] do |t|
	sh 'go', 'build', '-o', t.name, 'pihole/client'
end

file 'bin/piholed' => DEPS + FileList['src/pihole/**/*'] do |t|
	sh 'go', 'build', '-o', t.name, 'pihole/server'
end

task :default => ['bin/pihole', 'bin/piholed']

task :test do
	sh 'go', 'test', 'pihole/...'
end

task :clean do
	FileUtils.rm('bin/pihole')
	FileUtils.rm('bin/piholed')
end
