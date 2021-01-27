use survey;

CREATE TABLE IF NOT EXISTS `Catalogs` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `Title` varchar(200) DEFAULT NULL,
    `Description` varchar(1024) DEFAULT NULL,
    `Created` timestamp NULL DEFAULT NULL,
    `Updated` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`ID`),
    UNIQUE KEY `ID_UNIQUE` (`ID`)
);

CREATE TABLE IF NOT EXISTS `Questions` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `CatalogID` bigint(20) NOT NULL,
    `Question` varchar(200) DEFAULT NULL,
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`SurveyID`) REFERENCES `Catalogs`(`ID`)
);

CREATE TABLE IF NOT EXISTS `Options` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `QuestionID` bigint(20) NOT NULL,
    `Option` int NOT NULL,
    `Text` varchar(200) NOT NULL,
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`QuestionID`) REFERENCES `Questions`(`ID`)
);

CREATE TABLE IF NOT EXISTS `Answers` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `CatalogID` bigint(20) NOT NULL,
    `QuestionID` bigint(20) NOT NULL,
    `OptionID` bigint(20) NOT NULL,
    `SessionID` varchar(254) NOT NULL,
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`SurveyID`) REFERENCES `Catalogs`(`ID`),
    FOREIGN Key (`QuestionID`) REFERENCES `Questions`(`ID`),
    FOREIGN Key (`OptionID`) REFERENCES `Options`(`ID`)
);

CREATE TABLE `Users` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `Email` varchar(254) NOT NULL UNIQUE,
    `FirstName` varchar(254) NOT NULL,
    `LastName` varchar(254) NOT NULL,
    `Password` varchar(254) NOT NULL,
    `IsAdmin` boolean NOT NULL,
    `IsSuperuser` boolean NOT NULL,
    PRIMARY KEY (`ID`)
);

