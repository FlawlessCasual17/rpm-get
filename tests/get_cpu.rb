# frozen_string_literal: true

# Gets the CPU architecture in [String] format
# @return [String]
def retrieve_cpu
  if defined?(RUBY_PLATFORM)
    String(RUBY_PLATFORM).sub('-linux', '')
  elsif which('arch')
    `arch`.chomp
  elsif which('uname')
    `uname -m`.chomp
  end
end

HOST_CPU = retrieve_cpu
puts HOST_CPU
