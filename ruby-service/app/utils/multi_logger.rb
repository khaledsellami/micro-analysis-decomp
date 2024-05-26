require 'logger'
require 'singleton'
require 'fileutils'

class MultiLogger

  include Singleton

  PROGNAME = 'st_analyzer'
  def initialize
    @console_logger = Logger.new(STDOUT, progname: PROGNAME)
    @file_logger = nil
    set_default_logger_level
  end

  def set_default_logger_level
    @console_level = Logger::INFO
    @file_level = Logger::DEBUG
    @console_logger.level = @console_level
    @file_logger.level = @file_level if @file_logger!=nil
  end
  
  def set_logger_level(level: "default")
    if level == "default"
      @console_level = Logger::INFO
      @file_level = Logger::DEBUG
      @console_logger.level = @console_level
      @file_logger.level = @file_level if @file_logger!=nil
    else
      @console_logger.level = Logger.const_get(level.upcase)
      @file_logger.level = Logger.const_get(level.upcase) if @file_logger!=nil
    end
  end
  
  def add_file(file_path)
    FileUtils.mkdir_p(File.dirname(file_path))
    @file_logger = Logger.new(file_path, level: @file_level, progname: PROGNAME)
  end

  def debug(message)
    @console_logger.debug(message)
    @file_logger.debug(message) if @file_logger!=nil
  end

  def info(message)
    @console_logger.info(message)
    @file_logger.info(message) if @file_logger!=nil
  end

  def warn(message)
    @console_logger.warn(message)
    @file_logger.warn(message) if @file_logger!=nil
  end

  def error(message)
    @console_logger.error(message)
    @file_logger.error(message) if @file_logger!=nil
  end

  def fatal(message)
    @console_logger.fatal(message)
    @file_logger.fatal(message) if @file_logger!=nil
  end

  def unknown(message)
    @console_logger.unknown(message)
    @file_logger.unknown(message) if @file_logger!=nil
  end

  def log(message)
    @console_logger.log(message)
    @file_logger.log(message) if @file_logger!=nil
  end

  def close
    @file_logger.close if @file_logger!=nil
    @file_level = nil
    set_default_logger_level
  end
end