# frozen_string_literal: true

# Checks if a command is available
# @return [TrueClass, FalseClass]
def which(cmd)
  exts = ENV['PATHEXT'] ? ENV['PATHEXT'].split(File::PATH_SEPARATOR) : ['']
  ENV['PATH'].split(File::PATH_SEPARATOR).each do |path|
    exts.each do |ext|
      exe = File.join(path, "#{cmd}#{ext}")
      return true if File.executable?(exe) && !File.directory?(exe)
    end
  end
  false # Return false if not found
end

test = which('sudo')
puts test
