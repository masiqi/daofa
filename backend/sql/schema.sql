CREATE TABLE IF NOT EXISTS `admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `subject` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,  -- 可为空
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `knowledge_point` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `subject_id` int(11) NOT NULL,
  `parent_id` int(11) DEFAULT NULL,  -- 可为空，表示顶级知识点
  `name` varchar(255) NOT NULL,
  `description` text,  -- 可为空
  `is_leaf` boolean NOT NULL DEFAULT FALSE,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `subject_id` (`subject_id`),
  KEY `parent_id` (`parent_id`),
  KEY `idx_is_leaf` (`is_leaf`),
  CONSTRAINT `fk_knowledge_point_subject` FOREIGN KEY (`subject_id`) REFERENCES `subject` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_knowledge_point_parent` FOREIGN KEY (`parent_id`) REFERENCES `knowledge_point` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 删除不再使用的表
DROP TABLE IF EXISTS `exercise_question`;
DROP TABLE IF EXISTS `exercise_material`;

-- 创建题目类型表
CREATE TABLE IF NOT EXISTS `question_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,  -- 可为空
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建新的题目表
CREATE TABLE IF NOT EXISTS `question` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `image_path` varchar(255),  -- 可为空
  `ocr_text` text,  -- 可为空
  `answer` text NOT NULL,
  `explanation` text,  -- 可为空
  `type_id` int(11) NOT NULL,
  `hash` varchar(64) NOT NULL,  -- 新增 hash 字段
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hash` (`hash`),  -- 为 hash 字段添加唯一索引
  KEY `type_id` (`type_id`),
  CONSTRAINT `fk_question_type` FOREIGN KEY (`type_id`) REFERENCES `question_type` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建题目与知识点的关联表
CREATE TABLE IF NOT EXISTS `question_knowledge_point` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `question_id` int(11) NOT NULL,
  `knowledge_point_id` int(11) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_question_knowledge_point` (`question_id`, `knowledge_point_id`),
  KEY `question_id` (`question_id`),
  KEY `knowledge_point_id` (`knowledge_point_id`),
  CONSTRAINT `fk_qkp_question` FOREIGN KEY (`question_id`) REFERENCES `question` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_qkp_knowledge_point` FOREIGN KEY (`knowledge_point_id`) REFERENCES `knowledge_point` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;