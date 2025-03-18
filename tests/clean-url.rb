# typed: true
# frozen_string_literal: true

require 'faraday'

USER_AGENT = 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'

# Takes a url and cleans it of unnecessary data.
# @return [String]
def clean_url(url)
  # Create a Faraday connection that follows redirects
  conn = Faraday.new do |f|
    f.use Faraday::FollowRedirects::Middleware
  end

  # Make a HEAD request to get the final URL
  response = conn.head(url)

  # Get the final URL and trim any .rpm* patterns to just .rpm
  final_url = response.env.url.to_s
  final_url.sub(/\.rpm.*/, '.rpm')
end

test = clean_url('https://sourceforge.net/projects/openofficeorg.mirror/files/4.1.15/binaries/en-US/Apache_OpenOffice_4.1.15_Linux_x86-64_langpack-rpm_en-US.tar.gz/download')
puts test
