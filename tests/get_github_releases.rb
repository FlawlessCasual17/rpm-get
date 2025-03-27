# typed: true
# frozen_string_literal: true

require 'down'
require 'tty-progressbar'
require 'ruby-progressbar'

$CACHE_DIR = File.join(Dir.home, '.cache/rpm-get')
$USER_AGENT = 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36'

$header_auth = "Authorization: Bearer #{ENV['RPMGET_TOKEN']}"

# Creates the cache dir.
def create_cache_dir = Dir.mkdir($CACHE_DIR)

# Gets the package releases from GitHub
def get_github_releases
  url = 'https://api.github.com/repos/peazip/Peazip/releases/latest'
  c_file_path = "#{$CACHE_DIR}/test.json"

  progress_bar = nil # Initialize outside the block
  format = '[:bar] :percent TOTAL::total_byte :current/:total bytes ETA::eta :rate bytes/s'

  # Ensure cache directory exists
  create_cache_dir() unless Dir.exist?($CACHE_DIR)

  return if File.exist?(c_file_path)

  printc "Updating #{c_file_path}", 'info'
  printc "Downloading JSON cache of app to #{c_file_path}", 'progress'

  begin
    Down.download(
      url,
      destination: c_file_path,
      headers: {
        'User-Agent' => $USER_AGENT,
        'Authorization' => $header_auth
      },
      content_length_proc: lambda do |content_length|
        # If `content_length` is null, change the format
        content_length.nil? &&
          format = format.sub('TOTAL::total_byte', '').sub('/:total', '')

        bar = TTY::ProgressBar.new(format, total: content_length, head: '>')
        bar.resize(50)
        progress_bar = bar
      end,
      progress_proc: lambda do |progress|
        raise TypeError if progress.nil?

        progress_bar&.advance(progress - progress_bar.current)
      end
    )

    printc 'Download complete!', 'progress', true
  rescue StandardError => e
    printc "Failed to update #{c_file_path}", 'warn', true
    printc e.detailed_message, 'error'
  end
end

# Prints colored text to the terminal
def printc(msg, msg_type, new_line = false)
  # Define colours here
  @RED = "\e[31m"
  @GREEN = "\e[32m"
  @YELLOW = "\e[33m"
  @BLUE = "\e[34m"
  @MAGENTA = "\e[35m"
  @GREY = "\e[37m"
  @RESET = "\e[0m"

  # Add a new line if `new_line` is true
  cr = new_line ? "\n" : ''

  case msg_type
  when 'info'
    printf "#{cr}  [#{@GREEN}INFO#{@RESET}]: #{msg}\n"
  when 'progress'
    printf "#{cr}  [#{@BLUE}PROGRESS#{@RESET}]: #{msg}\n"
  when 'warn'
    printf "#{cr}  [#{@YELLOW}WARNING#{@RESET}]: #{msg}\n"
  when 'error'
    printf "#{cr}  [#{@RED}ERROR#{@RESET}]: #{msg}\n"
  when 'fatal'
    printf "#{cr}  [#{@MAGENTA}FATAL#{@RESET}]: #{msg}\n"
    exit 1
  else
    printf "#{cr}  [#{@GREY}UNKNOWN#{@RESET}]: #{msg}\n"
  end
end


# Testing
get_github_releases()
