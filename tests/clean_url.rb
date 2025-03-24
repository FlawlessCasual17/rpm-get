# typed: true
# frozen_string_literal: true

require 'faraday'
require 'typhoeus'

USER_AGENT = 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'

# Takes a url and cleans it of unnecessary data.
# @return [String]
def clean_url(url)
  # Create a HEAD request with follow_location enabled
  request = Typhoeus::Request.new(
    url,
    method: :head,
    followlocation: true,
    maxredirs: 10
  )

  # Execute the request
  response = request.run

  # Get the effective URL (final URL after following redirects)
  final_url = response.effective_url

  # Trim any .rpm* patterns to just .rpm
  String(final_url).sub(/\.rpm.*/, '.rpm')
end

test = clean_url('https://sourceforge.net/projects/openofficeorg.mirror/files/4.1.15/binaries/en-US/Apache_OpenOffice_4.1.15_Linux_x86-64_langpack-rpm_en-US.tar.gz/download')
puts test
