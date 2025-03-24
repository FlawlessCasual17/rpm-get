# typed: true
# frozen_string_literal: true

require 'fileutils'
require 'json'

# Retrieves JSON content from a file,
# then parses it into a Ruby object.
def parse_json(path, json_path)
  raise TypeError unless path.is_a?(String)
  raise TypeError unless json_path.is_a?(Array)

  # Read file and parse JSON content
  data = JSON.parse(File.read(path))
  # Access the nested data using the dig method
  # Using "&." for safe navigation
  # dig is a real method. See https://docs.ruby-lang.org/en/3.0/dig_methods_rdoc.html
  json_path.reduce(data) { |acc, key| acc&.dig(key) }
end

path = File.join(Dir.home, 'Downloads', 'releases-latest.json')
json_path = ['assets', 0]
value = parse_json(path, json_path)
puts value
