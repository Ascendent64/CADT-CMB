-- MySQL dump 10.13  Distrib 8.0.39, for Linux (x86_64)
--
-- Host: localhost    Database: CADTDatabase
-- ------------------------------------------------------
-- Server version	8.0.39-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `SequelizeMeta`
--

DROP TABLE IF EXISTS `SequelizeMeta`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `SequelizeMeta` (
  `name` varchar(255) NOT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `SequelizeMeta`
--

LOCK TABLES `SequelizeMeta` WRITE;
/*!40000 ALTER TABLE `SequelizeMeta` DISABLE KEYS */;
/*!40000 ALTER TABLE `SequelizeMeta` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `audit`
--

DROP TABLE IF EXISTS `audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `audit` (
  `id` int NOT NULL AUTO_INCREMENT,
  `orgUid` varchar(255) NOT NULL,
  `registryId` varchar(255) DEFAULT NULL,
  `rootHash` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `change` text,
  `table_name` varchar(255) DEFAULT NULL,
  `onchainConfirmationTimeStamp` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `audit`
--

LOCK TABLES `audit` WRITE;
/*!40000 ALTER TABLE `audit` DISABLE KEYS */;
/*!40000 ALTER TABLE `audit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cobenefits`
--

DROP TABLE IF EXISTS `cobenefits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cobenefits` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `cobenefit` text,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cobenefits`
--

LOCK TABLES `cobenefits` WRITE;
/*!40000 ALTER TABLE `cobenefits` DISABLE KEYS */;
INSERT INTO `cobenefits` VALUES ('73cfbe9c-8cea-4aca-94d8-f1641e686787','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','Sample Benefit',1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('74cfbe9c-8cea-4aca-94d8-f1641e686788','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','Sample Benefit 2',1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('75cfbe9c-8cea-4aca-94d8-f1641e686789','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','Sample Benefit 3',1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('85cfbe9c-8cea-4aca-94d8-f1641e686790','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','Sample Benefit 4',1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('96cfbe9c-8cea-4aca-94d8-f1641e686791','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','Sample Benefit 5',1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('a6cfbe9c-8cea-4aca-94d8-f1641e686792','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','Sample Benefit 6',1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `cobenefits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `estimations`
--

DROP TABLE IF EXISTS `estimations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `estimations` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `creditingPeriodStart` datetime DEFAULT NULL,
  `creditingPeriodEnd` datetime DEFAULT NULL,
  `unitCount` int DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `estimations`
--

LOCK TABLES `estimations` WRITE;
/*!40000 ALTER TABLE `estimations` DISABLE KEYS */;
INSERT INTO `estimations` VALUES ('c73fb4e7-3bd0-4449-8a57-6137b7c95a1f','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','2022-02-04 00:00:00','2022-03-04 00:00:00',100,1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('c74fb4e7-3bd0-4449-8a57-6137b7c95a2f','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','2022-02-05 00:00:00','2022-03-05 00:00:00',200,1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('c75fb4e7-3bd0-4449-8a57-6137b7c95a3f','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','2022-02-06 00:00:00','2022-03-06 00:00:00',300,1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('d75fb4e7-3bd0-4449-8a57-6137b7c95a4f','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','2022-03-06 00:00:00','2022-04-06 00:00:00',400,1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('e75fb4e7-3bd0-4449-8a57-6137b7c95a5f','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','2022-04-06 00:00:00','2022-05-06 00:00:00',500,1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('f75fb4e7-3bd0-4449-8a57-6137b7c95a6f','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','2022-05-06 00:00:00','2022-06-06 00:00:00',600,1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `estimations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `fileStore`
--

DROP TABLE IF EXISTS `fileStore`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `fileStore` (
  `id` varchar(255) NOT NULL,
  `orgUid` varchar(255) NOT NULL,
  `filePath` varchar(255) DEFAULT NULL,
  `fileType` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `fileStore`
--

LOCK TABLES `fileStore` WRITE;
/*!40000 ALTER TABLE `fileStore` DISABLE KEYS */;
/*!40000 ALTER TABLE `fileStore` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `governance`
--

DROP TABLE IF EXISTS `governance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `governance` (
  `id` varchar(255) NOT NULL,
  `orgUid` varchar(255) NOT NULL,
  `governanceType` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `governance`
--

LOCK TABLES `governance` WRITE;
/*!40000 ALTER TABLE `governance` DISABLE KEYS */;
/*!40000 ALTER TABLE `governance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `issuances`
--

DROP TABLE IF EXISTS `issuances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `issuances` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `startDate` datetime DEFAULT NULL,
  `endDate` datetime DEFAULT NULL,
  `verificationApproach` varchar(255) DEFAULT NULL,
  `verificationReportDate` datetime DEFAULT NULL,
  `verificationBody` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `issuances`
--

LOCK TABLES `issuances` WRITE;
/*!40000 ALTER TABLE `issuances` DISABLE KEYS */;
INSERT INTO `issuances` VALUES ('d9f58b08-af25-461c-88eb-403bb02b135e','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','2022-01-02 00:00:00','2022-02-11 00:00:00','Sample Approach','2022-03-16 00:00:00','Sample Body',1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('d9f68b08-af25-461c-88eb-403bb02b135f','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','2022-01-03 00:00:00','2022-02-12 00:00:00','Sample Approach 2','2022-03-17 00:00:00','Sample Body 2',1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('d9f78b08-af25-461c-88eb-403bb02b136f','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','2022-01-04 00:00:00','2022-02-13 00:00:00','Sample Approach 3','2022-03-18 00:00:00','Sample Body 3',1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('e9f78b08-af25-461c-88eb-403bb02b137f','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','2022-02-04 00:00:00','2022-03-13 00:00:00','Sample Approach 4','2022-04-18 00:00:00','Sample Body 4',1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('f9f78b08-af25-461c-88eb-403bb02b138f','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','2022-03-04 00:00:00','2022-04-13 00:00:00','Sample Approach 5','2022-05-18 00:00:00','Sample Body 5',1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('g9f78b08-af25-461c-88eb-403bb02b139f','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','2022-04-04 00:00:00','2022-05-13 00:00:00','Sample Approach 6','2022-06-18 00:00:00','Sample Body 6',1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `issuances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `label_unit`
--

DROP TABLE IF EXISTS `label_unit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `label_unit` (
  `id` varchar(255) NOT NULL,
  `warehouseUnitId` varchar(255) DEFAULT NULL,
  `label` varchar(255) DEFAULT NULL,
  `unitQuantity` int DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `label_unit`
--

LOCK TABLES `label_unit` WRITE;
/*!40000 ALTER TABLE `label_unit` DISABLE KEYS */;
/*!40000 ALTER TABLE `label_unit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `labels`
--

DROP TABLE IF EXISTS `labels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `labels` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `label` varchar(255) DEFAULT NULL,
  `labelType` varchar(255) DEFAULT NULL,
  `creditingPeriodStartDate` datetime DEFAULT NULL,
  `creditingPeriodEndDate` datetime DEFAULT NULL,
  `validityPeriodStartDate` datetime DEFAULT NULL,
  `validityPeriodEndDate` datetime DEFAULT NULL,
  `unitQuantity` int DEFAULT NULL,
  `labelLink` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `labels`
--

LOCK TABLES `labels` WRITE;
/*!40000 ALTER TABLE `labels` DISABLE KEYS */;
/*!40000 ALTER TABLE `labels` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `meta`
--

DROP TABLE IF EXISTS `meta`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `meta` (
  `id` varchar(255) NOT NULL,
  `metaKey` varchar(255) DEFAULT NULL,
  `metaValue` text,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `meta`
--

LOCK TABLES `meta` WRITE;
/*!40000 ALTER TABLE `meta` DISABLE KEYS */;
/*!40000 ALTER TABLE `meta` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `organizations`
--

DROP TABLE IF EXISTS `organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `organizations` (
  `orgUid` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `icon` varchar(255) DEFAULT NULL,
  `isHome` tinyint(1) DEFAULT NULL,
  `subscribed` tinyint(1) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`orgUid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `organizations`
--

LOCK TABLES `organizations` WRITE;
/*!40000 ALTER TABLE `organizations` DISABLE KEYS */;
/*!40000 ALTER TABLE `organizations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projectlocations`
--

DROP TABLE IF EXISTS `projectlocations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projectlocations` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `country` varchar(255) DEFAULT NULL,
  `inCountryRegion` varchar(255) DEFAULT NULL,
  `geographicIdentifier` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projectlocations`
--

LOCK TABLES `projectlocations` WRITE;
/*!40000 ALTER TABLE `projectlocations` DISABLE KEYS */;
INSERT INTO `projectlocations` VALUES ('8182100d-7794-4df7-b3b3-758391d13011','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','Latvia',NULL,'Sample Identifier',1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('8282100d-7794-4df7-b3b3-758391d13012','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','Estonia',NULL,'Sample Identifier 2',1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('8382100d-7794-4df7-b3b3-758391d13013','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','Lithuania',NULL,'Sample Identifier 3',1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('9482100d-7794-4df7-b3b3-758391d13014','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','Australia',NULL,'Sample Identifier 4',1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('a482100d-7794-4df7-b3b3-758391d13015','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','Brazil',NULL,'Sample Identifier 5',1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('b482100d-7794-4df7-b3b3-758391d13016','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','Indonesia',NULL,'Sample Identifier 6',1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `projectlocations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projectratings`
--

DROP TABLE IF EXISTS `projectratings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projectratings` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `ratingType` varchar(255) DEFAULT NULL,
  `ratingRangeHighest` varchar(255) DEFAULT NULL,
  `ratingRangeLowest` varchar(255) DEFAULT NULL,
  `rating` varchar(255) DEFAULT NULL,
  `ratingLink` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projectratings`
--

LOCK TABLES `projectratings` WRITE;
/*!40000 ALTER TABLE `projectratings` DISABLE KEYS */;
INSERT INTO `projectratings` VALUES ('d31c3c75-b944-498d-9557-315f9005f478','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','CCQI','100','0','97','testlink.com',1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('d41c3c75-b944-498d-9557-315f9005f479','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','CCQI','100','0','98','testlink2.com',1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('d51c3c75-b944-498d-9557-315f9005f480','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','CCQI','100','0','99','testlink3.com',1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('e51c3c75-b944-498d-9557-315f9005f481','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','CCQI','100','0','95','testlink4.com',1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('f51c3c75-b944-498d-9557-315f9005f482','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','CCQI','100','0','96','testlink5.com',1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('g51c3c75-b944-498d-9557-315f9005f483','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','CCQI','100','0','97','testlink6.com',1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `projectratings` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `projects`
--

DROP TABLE IF EXISTS `projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `projects` (
  `warehouseProjectId` varchar(255) NOT NULL,
  `orgUid` varchar(255) NOT NULL,
  `currentRegistry` varchar(255) DEFAULT NULL,
  `projectId` varchar(255) DEFAULT NULL,
  `originProjectId` varchar(255) DEFAULT NULL,
  `registryOfOrigin` varchar(255) DEFAULT NULL,
  `program` varchar(255) DEFAULT NULL,
  `projectName` varchar(255) NOT NULL,
  `projectLink` varchar(255) DEFAULT NULL,
  `projectDeveloper` varchar(255) DEFAULT NULL,
  `sector` varchar(255) DEFAULT NULL,
  `projectType` varchar(255) DEFAULT NULL,
  `projectTags` text,
  `coveredByNDC` varchar(255) DEFAULT NULL,
  `ndcInformation` text,
  `projectStatus` varchar(255) DEFAULT NULL,
  `projectStatusDate` datetime DEFAULT NULL,
  `unitMetric` varchar(255) DEFAULT NULL,
  `methodology` text,
  `validationBody` varchar(255) DEFAULT NULL,
  `validationDate` datetime DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`warehouseProjectId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `projects`
--

LOCK TABLES `projects` WRITE;
/*!40000 ALTER TABLE `projects` DISABLE KEYS */;
INSERT INTO `projects` VALUES ('2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','Verra','790','124','Sweden National Registry',NULL,'Prevent Deforestation','deforestationtest.com','Dev 3','Fugitive emissions from fuel (solid, oil and gas)','Deforestation Prevention',NULL,'Outside NDC',NULL,'Registered','2022-03-02 00:00:00','tCO2e','Substitution of CO2 from fossil or mineral origin by CO2 from biogenic residual sources in the production of organic compounds --- Version 4.0',NULL,NULL,1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','EcoRegistry','791','125','Sweden National Registry',NULL,'Reduce Emissions','emissiontest.com','Dev 4','Fugitive emissions from fuel (solid, oil and gas)','Emission Reduction',NULL,'Outside NDC',NULL,'Registered','2022-04-02 00:00:00','tCO2e','Substitution of CO2 from fossil or mineral origin by CO2 from biogenic residual sources in the production of chemical compounds --- Version 5.0',NULL,NULL,1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','BioCarbon Standard','792','126','Tero Carbon',NULL,'Forest Conservation','forestconservationtest.com','Dev 5','Agriculture; forestry and fishing','CAR - Forest',NULL,'Inside NDC',NULL,'Registered','2022-05-02 00:00:00','tCO2e','VCS - VM0027',NULL,NULL,1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','Royal Kingdom of Bhutan','793','127','Royal Kingdom of Bhutan',NULL,'Soil Restoration','soilrestorationtest.com','Dev 6','Agriculture; forestry and fishing','GS - Soil Organic Carbon Activity Module for Application of Organic Soil Improvers from Pulp and Paper Mill Sludges',NULL,'Inside NDC',NULL,'Registered','2022-06-02 00:00:00','tCO2e','VCS - VM0041',NULL,NULL,1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','CDM Registry','794','128','CDM Registry',NULL,'Methane Destruction','methanedestructiontest.com','Dev 7','Carbon capture and storage','ACR - Truck Stop Electrification',NULL,'Inside NDC',NULL,'Registered','2022-07-02 00:00:00','tCO2e','CDM - AM0026',NULL,NULL,1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55'),('9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','Global Carbon Council','789','123','Sweden National Registry',NULL,'Stop Desertification','desertificationtest.com','Dev 2','Fugitive emissions from fuel (solid, oil and gas)','Coal Mine Methane',NULL,'Outside NDC',NULL,'Registered','2022-02-02 00:00:00','tCO2e','Substitution of CO2 from fossil or mineral origin by CO2 from biogenic residual sources in the production of inorganic compounds --- Version 3.0',NULL,NULL,1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55');
/*!40000 ALTER TABLE `projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `relatedprojects`
--

DROP TABLE IF EXISTS `relatedprojects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `relatedprojects` (
  `id` varchar(255) NOT NULL,
  `warehouseProjectId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `relatedProjectId` varchar(255) DEFAULT NULL,
  `relationshipType` varchar(255) DEFAULT NULL,
  `registry` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `relatedprojects`
--

LOCK TABLES `relatedprojects` WRITE;
/*!40000 ALTER TABLE `relatedprojects` DISABLE KEYS */;
INSERT INTO `relatedprojects` VALUES ('e880047e-cdf4-45bb-a9df-e706fa427713','9b9bb857-c71b-4649-b805-a289db27dc1c','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','333','Sample',NULL,1646975765,'2022-03-11 05:17:55','2022-03-11 05:17:55'),('e890047e-cdf4-45bb-a9df-e706fa427714','2d9bb857-c71b-4649-b805-a289db27dc1d','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','334','Sample 2',NULL,1646975765,'2022-03-12 05:17:55','2022-03-12 05:17:55'),('e910047e-cdf4-45bb-a9df-e706fa427715','3e9bb857-c71b-4649-b805-a289db27dc1e','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','335','Sample 3',NULL,1646975765,'2022-03-13 05:17:55','2022-03-13 05:17:55'),('f910047e-cdf4-45bb-a9df-e706fa427716','4f9bb857-c71b-4649-b805-a289db27dc1f','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','336','Sample 4',NULL,1646975765,'2022-04-11 04:17:55','2022-04-11 04:17:55'),('g910047e-cdf4-45bb-a9df-e706fa427717','5g9bb857-c71b-4649-b805-a289db27dc1g','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','337','Sample 5',NULL,1646975765,'2022-05-11 04:17:55','2022-05-11 04:17:55'),('h910047e-cdf4-45bb-a9df-e706fa427718','6h9bb857-c71b-4649-b805-a289db27dc1h','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','338','Sample 6',NULL,1646975765,'2022-06-11 04:17:55','2022-06-11 04:17:55');
/*!40000 ALTER TABLE `relatedprojects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `simulator`
--

DROP TABLE IF EXISTS `simulator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `simulator` (
  `id` varchar(255) NOT NULL,
  `orgUid` varchar(255) NOT NULL,
  `simulatorType` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `simulator`
--

LOCK TABLES `simulator` WRITE;
/*!40000 ALTER TABLE `simulator` DISABLE KEYS */;
/*!40000 ALTER TABLE `simulator` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `staging`
--

DROP TABLE IF EXISTS `staging`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `staging` (
  `id` int NOT NULL AUTO_INCREMENT,
  `table_name` varchar(255) DEFAULT NULL,
  `action` varchar(255) DEFAULT NULL,
  `data` text,
  `commited` tinyint(1) DEFAULT NULL,
  `failedCommit` tinyint(1) DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `staging`
--

LOCK TABLES `staging` WRITE;
/*!40000 ALTER TABLE `staging` DISABLE KEYS */;
/*!40000 ALTER TABLE `staging` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `units`
--

DROP TABLE IF EXISTS `units`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `units` (
  `unitBlockStart` varchar(255) DEFAULT NULL,
  `unitBlockEnd` varchar(255) DEFAULT NULL,
  `unitCount` int DEFAULT NULL,
  `warehouseUnitId` varchar(255) NOT NULL,
  `issuanceId` varchar(255) DEFAULT NULL,
  `projectLocationId` varchar(255) DEFAULT NULL,
  `orgUid` varchar(255) NOT NULL,
  `unitOwner` varchar(255) DEFAULT NULL,
  `countryJurisdictionOfOwner` varchar(255) DEFAULT NULL,
  `inCountryJurisdictionOfOwner` varchar(255) DEFAULT NULL,
  `serialNumberBlock` varchar(255) DEFAULT NULL,
  `serialNumberPattern` varchar(255) DEFAULT NULL,
  `vintageYear` int DEFAULT NULL,
  `unitType` varchar(255) DEFAULT NULL,
  `marketplace` varchar(255) DEFAULT NULL,
  `marketplaceLink` varchar(255) DEFAULT NULL,
  `marketplaceIdentifier` varchar(255) DEFAULT NULL,
  `unitTags` text,
  `unitStatus` varchar(255) DEFAULT NULL,
  `unitStatusReason` text,
  `unitRegistryLink` varchar(255) DEFAULT NULL,
  `correspondingAdjustmentDeclaration` varchar(255) DEFAULT NULL,
  `correspondingAdjustmentStatus` varchar(255) DEFAULT NULL,
  `timeStaged` bigint DEFAULT NULL,
  `createdAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`warehouseUnitId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `units`
--

LOCK TABLES `units` WRITE;
/*!40000 ALTER TABLE `units` DISABLE KEYS */;
INSERT INTO `units` VALUES ('E345','F567',666,'09d7a102-a5a6-4f80-bc67-d28eba4952f6','d9f78b08-af25-461c-88eb-403bb02b136f','8382100d-7794-4df7-b3b3-758391d13013','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','Sample Owner 3','Lithuania',NULL,'E345-F567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2018,'Reduction - technical 3',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl3.com','Unknown','Pending',1647141471,'2022-03-15 05:29:39','2022-03-15 05:29:39'),('G345','H567',888,'19d7a102-a5a6-4f80-bc67-d28eba4952f7','e9f78b08-af25-461c-88eb-403bb02b137f','9482100d-7794-4df7-b3b3-758391d13014','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','Sample Owner 4','Australia',NULL,'G345-H567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2020,'Reduction - technical 4',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl4.com','Unknown','Pending',1647141471,'2022-04-15 04:29:39','2024-09-06 03:36:46'),('B568','C789',333,'19d8a103-b5b7-5g91-bd68-e29fcb5063f4','d9f58b08-af25-461c-88eb-403bb02b135e','8182100d-7794-4df7-b3b3-758391d13011','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','Another Owner','Belize',NULL,'B568-C789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2015,'Reduction - technical',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl.com','Unknown','Pending',1647141472,'2022-03-13 05:29:39','2024-09-06 03:36:46'),('D568','E789',555,'29d8a103-b5b7-5g91-bd68-e29fcb5063f5','d9f68b08-af25-461c-88eb-403bb02b135f','8282100d-7794-4df7-b3b3-758391d13012','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','Another Owner 2','Estonia',NULL,'D568-E789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2017,'Reduction - technical 2',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl2.com','Unknown','Pending',1647141472,'2022-03-14 05:29:39','2024-09-06 03:36:46'),('H568','I789',999,'29d8a103-b5b7-5g91-bd68-e29fcb5063f8','e9f78b08-af25-461c-88eb-403bb02b137f','9482100d-7794-4df7-b3b3-758391d13014','97741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d7','Another Owner 4','Australia',NULL,'H568-I789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2021,'Reduction - technical 4',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl4.com','Unknown','Pending',1647141472,'2022-04-15 04:29:39','2024-09-06 03:36:46'),('I345','J567',1111,'39d7a102-a5a6-4f80-bc67-d28eba4952f9','f9f78b08-af25-461c-88eb-403bb02b138f','a482100d-7794-4df7-b3b3-758391d13015','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','Sample Owner 5','Brazil',NULL,'I345-J567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2022,'Reduction - technical 5',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl5.com','Unknown','Pending',1647141471,'2022-05-15 04:29:39','2024-09-06 03:36:46'),('F568','G789',777,'39d8a103-b5b7-5g91-bd68-e29fcb5063f6','d9f78b08-af25-461c-88eb-403bb02b136f','8382100d-7794-4df7-b3b3-758391d13013','87641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2e8','Another Owner 3','Lithuania',NULL,'F568-G789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2019,'Reduction - technical 3',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl3.com','Unknown','Pending',1647141472,'2022-03-15 05:29:39','2024-09-06 03:36:46'),('J568','K789',1222,'49d8a103-b5b7-5g91-bd68-e29fcb5063f1','f9f78b08-af25-461c-88eb-403bb02b138f','a482100d-7794-4df7-b3b3-758391d13015','a7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d6','Another Owner 5','Brazil',NULL,'J568-K789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2023,'Reduction - technical 5',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl5.com','Unknown','Pending',1647141472,'2022-05-15 04:29:39','2024-09-06 03:36:46'),('K345','L567',1333,'59d7a102-a5a6-4f80-bc67-d28eba4952f1','g9f78b08-af25-461c-88eb-403bb02b139f','b482100d-7794-4df7-b3b3-758391d13016','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','Sample Owner 6','Indonesia',NULL,'K345-L567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2024,'Reduction - technical 6',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl6.com','Unknown','Pending',1647141471,'2022-06-15 04:29:39','2024-09-06 03:36:46'),('L568','M789',1444,'69d8a103-b5b7-5g91-bd68-e29fcb5063f1','g9f78b08-af25-461c-88eb-403bb02b139f','b482100d-7794-4df7-b3b3-758391d13016','b7741db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d5','Another Owner 6','Indonesia',NULL,'L568-M789','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2025,'Reduction - technical 6',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl6.com','Unknown','Pending',1647141472,'2022-06-15 04:29:39','2024-09-06 03:36:46'),('A345','B567',222,'89d7a102-a5a6-4f80-bc67-d28eba4952f3','d9f58b08-af25-461c-88eb-403bb02b135e','8182100d-7794-4df7-b3b3-758391d13011','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d9','Sample Owner','Belize',NULL,'A345-B567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2014,'Reduction - technical',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl.com','Unknown','Pending',1647141471,'2022-03-13 05:29:39','2022-03-13 05:29:39'),('C345','D567',444,'99d7a102-a5a6-4f80-bc67-d28eba4952f4','d9f68b08-af25-461c-88eb-403bb02b135f','8282100d-7794-4df7-b3b3-758391d13012','77641db780adc6c74f1ff357804e26a799e4a09157f426aac588963a39bdb2d8','Sample Owner 2','Estonia',NULL,'C345-D567','[.*\\D]+([0-9]+)+[-][.*\\D]+([0-9]+)$',2016,'Reduction - technical 2',NULL,NULL,'eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT',NULL,'Buffer',NULL,'sampleurl2.com','Unknown','Pending',1647141471,'2022-03-14 05:29:39','2022-03-14 05:29:39');
/*!40000 ALTER TABLE `units` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-09-08 20:10:45
