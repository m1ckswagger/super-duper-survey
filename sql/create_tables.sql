use survey;

CREATE TABLE IF NOT EXISTS `Catalogs` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `Title` varchar(200) DEFAULT NULL,
    `Description` varchar(1024) DEFAULT NULL,
    `Created` timestamp NULL DEFAULT NULL,
    `Updated` timestamp NULL DEFAULT NULL,
    `Due` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`ID`),
    UNIQUE KEY `ID_UNIQUE` (`ID`)
);

CREATE TABLE IF NOT EXISTS `Questions` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `CatalogID` bigint(20) NOT NULL,
    `Question` varchar(8000) DEFAULT NULL,
    `Num` bigint(20) NOT NULL,
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`CatalogID`) REFERENCES `Catalogs`(`ID`)
);

CREATE TABLE IF NOT EXISTS `Options` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `QuestionID` bigint(20) NOT NULL,
    `Num` bigint(20) NOT NULL,
    `Text` varchar(200) NOT NULL,
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`QuestionID`) REFERENCES `Questions`(`ID`)
);

CREATE TABLE IF NOT EXISTS `Answers` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `CatalogID` bigint(20) NOT NULL,
    `QuestionNum` bigint(20) NOT NULL,
    `OptionNum` bigint(20) NOT NULL,
    `SessionID` varchar(254) NOT NULL,
    PRIMARY KEY (`ID`)
/*
    FOREIGN KEY (`CatalogID`) REFERENCES `Catalogs`(`ID`),
    FOREIGN Key (`QuestionNum`) REFERENCES `Questions`(`Num`),
    FOREIGN Key (`OptionNum`) REFERENCES `Options`(`Num`)
*/
);

CREATE TABLE IF NOT EXISTS `Users` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `Email` varchar(254) NOT NULL UNIQUE,
    `FirstName` varchar(254) NOT NULL,
    `LastName` varchar(254) NOT NULL,
    `Password` varchar(254) NOT NULL,
    `IsAdmin` boolean NOT NULL,
    `IsSuperuser` boolean NOT NULL,
    PRIMARY KEY (`ID`)
);

