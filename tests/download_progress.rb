# typed: true
# frozen_string_literal: true

require 'down'
require 'tty-progressbar'

# Define the URL and the destination file path
url = 'http://ipv4.download.thinkbroadband.com/5MB.zip'
file_name = url.split('/')[-1]
file_path = "/home/kris/.cache/#{file_name}"

# Initialize the progress bar
progress_bar = nil

# formatting
format_str = "    Downloading \"#{file_name}\" [:bar] :percent TOTAL::total_byte :current/:total bytes ETA::eta :rate/s   "

# Download the file with progress tracking
Down.download(
  url,
  destination: file_path,
  content_length_proc: lambda do |content_length|
    raise TypeError if content_length.nil?

    bar = TTY::ProgressBar.new(format_str, total: content_length)
    bar.resize(30)
    progress_bar = bar
  end,
  progress_proc: lambda do |progress|
    raise TypeError if progress.nil?

    progress_bar&.advance(progress - progress_bar.current)
  end
)

puts "\nDownload completed! File saved to #{file_path}"
