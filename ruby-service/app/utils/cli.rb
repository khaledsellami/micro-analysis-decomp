require 'optparse'
require 'ostruct'

class Cli
  def self.parse(options)
    args = OpenStruct.new
    args.logging = 'default'
    args.output = File.join(Dir.pwd, 'data', 'ruby')
    args.monolithic = false

    opt_parser = OptionParser.new do |opts|
      opts.banner = 'Usage: st_analyzer [options]'

      opts.on('-p', '--path PATH', 'The path to source code of the application') do |path|
        args.path = path
      end

      opts.on('-o', '--output OUTPUT', 'The output path to save the results in') do |output|
        args.output = output
      end

      # add logging options handling
      opts.on('-l', '--logging LOGGING', 'The logging level', %w[default info debug warning error]) { |logging|
        args.logging = logging }

      opts.on('-m', '--monolithic', 'To specify is the application being analyzed is monolithic or not.') do |monolithic|
        args.monolithic = monolithic
      end

      opts.on('-h', '--help', 'Prints this help') do
        puts opts
        exit
      end
    end

    opt_parser.parse!(options)
    raise OptionParser::MissingArgument if args.path.nil?

    args
  end
end
