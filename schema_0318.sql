-- MySQL Script generated by MySQL Workbench
-- Fri Mar 18 12:47:45 2022
-- Model: New Model    Version: 1.0
-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema OTTsdb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema OTTsdb
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `OTTsdb`;
CREATE SCHEMA IF NOT EXISTS `OTTsdb` DEFAULT CHARACTER SET utf8 ;
USE `OTTsdb` ;

-- -----------------------------------------------------
-- Table `OTTsdb`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `OTTsdb`.`user` (
  `userId` INT(64) NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NULL,
  `accessToken` VARCHAR(400) NULL,
  PRIMARY KEY (`userId`),
  UNIQUE INDEX `userNo_UNIQUE` (`userId` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `OTTsdb`.`OTTServices`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `OTTsdb`.`OTTServices` (
  `OTTservicesId` INT(64) NOT NULL AUTO_INCREMENT,
  `platform` VARCHAR(45) NOT NULL,
  `membership` VARCHAR(45) NOT NULL,
  `price` INT(64) NOT NULL,
  PRIMARY KEY (`OTTservicesId`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `OTTsdb`.`subscribedServices`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `OTTsdb`.`subscribedServices` (
  `subscribedServiceId` INT(64) NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(45) NOT NULL,
  `OTTServiceId` INT(64) NOT NULL,
  `subscribedDate` DATE NOT NULL,
  `expireDate` DATE NOT NULL,
  `paymentType` INT(64) NOT NULL,
  PRIMARY KEY (`subscribedServiceId`),
  INDEX `fk_subscribedService_OTTServiceId_idx` (`OTTServiceId` ASC) VISIBLE,
  CONSTRAINT `fk_subscribedService_OTTServiceId`
    FOREIGN KEY (`OTTServiceId`)
    REFERENCES `OTTsdb`.`OTTServices` (`OTTservicesId`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
