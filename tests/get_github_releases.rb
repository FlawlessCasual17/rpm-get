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
  cache_file = "#{$CACHE_DIR}/test.json"

  # Ensure cache directory exists
  Dir.mkdir($CACHE_DIR) unless Dir.exist?($CACHE_DIR)

  if !File.exist?(cache_file)
    printc "Updating #{cache_file}", 'info'

    total_size = 0
    progressbar = nil  # Initialize outside the block

    begin
      Down.download(
        'https://api.github.com/repos/peazip/Peazip/releases/latest',
        destination: cache_file,
        headers: {
          'User-Agent' => $USER_AGENT,
          'Authorization' => $header_auth
        },
        content_length_proc: lambda do |content_length|
          if content_length
            total_size = content_length
            # Initialize progress bar with total file size
            progressbar = TTY::ProgressBar.new(
              'Downloading [:bar] :percent :eta',
              total: total_size,
              width: 40
            )
          end
        end,
        progress_proc: lambda do |progress|
          progressbar&.advance(progress - progressbar.current)
        end
      )
    rescue StandardError => e
      printc "Failed to update #{cache_file}", 'warn'
      printc "#{e}\n#{e.detailed_message}", 'error'
    end
  end
end

# Prints colored text to the terminal
def printc(msg, msg_type)
  # Define colours here
  @RED = "\e[31m"
  @GREEN = "\e[32m"
  @YELLOW = "\e[33m"
  @BLUE = "\e[34m"
  @MAGENTA = "\e[35m"
  @GREY = "\e[37m"
  @RESET = "\e[0m"

  case msg_type
  when 'info'
    printf "  [#{@GREEN}INFO#{@RESET}]: #{msg}\n"
  when 'progress'
    printf "  [#{@BLUE}PROGRESS#{@RESET}]: #{msg}\n"
  when 'warn'
    printf "  [#{@YELLOW}WARNING#{@RESET}]: #{msg}\n"
  when 'error'
    printf "  [#{@RED}ERROR#{@RESET}]: #{msg}\n"
  when 'fatal'
    printf "  [#{@MAGENTA}FATAL#{@RESET}]: #{msg}\n"
    exit 1
  else
    printf "  [#{@GREY}UNKNOWN#{@RESET}]: #{msg}\n"
  end
end


# Testing
get_github_releases()
