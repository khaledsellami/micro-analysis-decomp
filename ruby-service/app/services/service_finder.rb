# frozen_string_literal: true

require_relative '../utils/multi_logger'

class ServiceFinder
  def initialize(app_path)
    @app_path = app_path
    @logger = MultiLogger.instance
  end

  def find_services
    # code here
    # TODO: Implement the service finder
    @logger.error('ServiceFinder#find_services not implemented yet! Returning the whole application as a single service.')
    { File.basename(@app_path) => @app_path}
  end
end
