# frozen_string_literal: true
# typed: true

# FOR LATER: @param [Array] xml_path

require 'nokogiri'
require 'rest-client'

# Parses HTML/XML content from a URL.
# @param [String] url
def parse_website(url)
  # @type [Object]
  content = RestClient.get(url)
  # @type [Nokogiri::HTML::Document]
  page = Nokogiri::HTML(content)

  puts page.content
end

# testing
parse_website('http://en.wikipedia.org/')
