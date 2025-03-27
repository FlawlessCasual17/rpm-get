# typed: true
# frozen_string_literal: true

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

puts printc('a message', 'info')
