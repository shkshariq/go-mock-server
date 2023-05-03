DROP TABLE IF EXISTS `mock_apis`;
CREATE TABLE `mock_apis` (
  `id` int NOT NULL AUTO_INCREMENT,
  `api_path` varchar(255) DEFAULT NULL,
  `status_code` int DEFAULT NULL,
  `headers` text,
  `body` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `api_path` (`api_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;