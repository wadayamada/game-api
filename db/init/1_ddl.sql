-- MySQL Script generated by MySQL Workbench
-- Fri Feb 14 23:09:20 2020
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema dojo_api
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema dojo_api
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `dojo_api` DEFAULT CHARACTER SET utf8mb4 ;
USE `dojo_api` ;

SET CHARSET utf8mb4;

-- -----------------------------------------------------
-- Table `dojo_api`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dojo_api`.`user` (
  `id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
  `auth_token` VARCHAR(128) NOT NULL COMMENT '認証トークン',
  `name` VARCHAR(64) NOT NULL COMMENT 'ユーザ名',
  `high_score` INT UNSIGNED NOT NULL COMMENT 'ハイスコア',
  `coin` INT UNSIGNED NOT NULL COMMENT '所持コイン',
  PRIMARY KEY (`id`),
  INDEX `idx_auth_token` (`auth_token` ASC))
ENGINE = InnoDB
COMMENT = 'ユーザ';


-- -----------------------------------------------------
-- Table `dojo_api`.`collection_item`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dojo_api`.`collection_item` (
  `id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
  `name` VARCHAR(64) NOT NULL COMMENT 'コレクションアイテム名',
  `rarity` INT NOT NULL COMMENT 'レアリティ',
  PRIMARY KEY (`id`))
ENGINE = InnoDB
COMMENT = 'コレクションアイテム';


-- -----------------------------------------------------
-- Table `dojo_api`.`user_collection_item`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dojo_api`.`user_collection_item` (
  `user_id` VARCHAR(128) NOT NULL COMMENT 'ユーザID',
  `collection_item_id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
  PRIMARY KEY (`user_id`, `collection_item_id`),
  INDEX `fk_user_collection_item_user_idx` (`user_id` ASC),
  INDEX `fk_user_collection_item_collection_item_idx` (`collection_item_id` ASC),
  CONSTRAINT `fk_user_collection_item_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `dojo_api`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_collection_item_collection_item`
    FOREIGN KEY (`collection_item_id`)
    REFERENCES `dojo_api`.`collection_item` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'ユーザ所持コレクションアイテム';


-- -----------------------------------------------------
-- Table `dojo_api`.`gacha_probability`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `dojo_api`.`gacha_probability` (
  `collection_item_id` VARCHAR(128) NOT NULL COMMENT 'コレクションアイテムID',
  `ratio` INT UNSIGNED NOT NULL COMMENT '排出重み',
  INDEX `fk_gacha_probability_collection_item_idx` (`collection_item_id` ASC),
  PRIMARY KEY (`collection_item_id`),
  CONSTRAINT `fk_gacha_probability_collection_item_id`
    FOREIGN KEY (`collection_item_id`)
    REFERENCES `dojo_api`.`collection_item` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'ガチャ排出情報';


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
