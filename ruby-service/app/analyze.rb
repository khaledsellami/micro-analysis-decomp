# frozen_string_literal: true

require 'fileutils'
require 'json'

require_relative 'utils/multi_logger'
require_relative 'parsers/project_parser'
require_relative 'services/service_finder'


def analyze_app(app_name, app_path, output_path: File.join(Dir.pwd, 'data', 'ruby'), logging_level: "default",
                is_monolithic: false)
  # Set up the logger
  logging_file = File.join(output_path, app_name, 'log.txt')
  logger = MultiLogger.instance
  logger.add_file(logging_file)
  logger.set_logger_level(level: logging_level)

  # find the roots of the different microservices
  logger.info("Processing application: #{app_name}")
  service_map = {}
  if is_monolithic
    logger.debug("Application is monolithic. Using root directory as the only service.")
    service_map[app_name] = app_path
  else
    logger.debug("Searching for microservices.")
    service_finder = ServiceFinder.new(app_path)
    service_map = service_finder.find_services
  end
  logger.debug("Found #{service_map.length} services.")

  # analyze each microservice
  all_objects = []
  all_executables = []
  service_map.each do |service_name, service_path|
    logger.debug("Working on service: #{service_name}")
    ast_processor = ProjectParser.new(service_path, service_name)
    ast_processor.parse
    all_objects.concat(ast_processor.objects_)
    all_executables.concat(ast_processor.executables_)
  end
  logger.info("Detected #{all_objects.length} classes")
  logger.info("Detected #{all_executables.length} methods")
  logger.info("Detected #{service_map.length} microservices")

  # save results
  save_path = File.join(output_path, app_name)
  logger.info("Saving data to: #{save_path}")
  save_data(all_objects, all_executables, save_path)

  # close the logger
  logger.close
end


def save_data(all_objects, all_executables, save_path)
  logger = MultiLogger.instance
  FileUtils.mkdir_p(save_path)
  objects_path = File.join(save_path, 'typeData.json')
  objects_hashes = all_objects.map(&:to_hash)
  objects_json = JSON.pretty_generate(objects_hashes)
  logger.debug("Saving class data in #{objects_path}")
  File.write(objects_path, objects_json)
  executables_path = File.join(save_path, 'methodData.json')
  executables_hashes = all_executables.map(&:to_hash)
  executables_json = JSON.pretty_generate(executables_hashes)
  logger.debug("Saving method data in #{executables_path}")
  File.write(executables_path, executables_json)
end

