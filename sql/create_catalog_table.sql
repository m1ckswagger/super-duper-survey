use survey;

CREATE TABLE `Catalogs` (
    `ID` bigint(20) NOT NULL AUTO_INCREMENT,
    `Title` varchar(200) DEFAULT NULL,
    `Description` varchar(1024) DEFAULT NULL,
    `Created` timestamp NULL DEFAULT NULL,
    `Updated` timestamp NULL DEFAULT NULL,
    PRIMARY KEY (`ID`),
    UNIQUE KEY `ID_UNIQUE` (`ID`)
);