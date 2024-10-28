CREATE TABLE agents
(
    id             INT AUTO_INCREMENT
        PRIMARY KEY,
    name           VARCHAR(20)   NOT NULL,
    number         VARCHAR(15)   NULL,
    email          VARCHAR(100)  NULL,
    hashedPassword VARCHAR(255)  NOT NULL,
    priority       INT DEFAULT 0 NOT NULL,
    isAdmin        TINYINT(1)    NOT NULL
);

CREATE TABLE blockedPhoneNumbers
(
    id          INT AUTO_INCREMENT
        PRIMARY KEY,
    phoneNumber VARCHAR(15) NOT NULL,
    blockDate   DATETIME    NOT NULL
);

CREATE TABLE calls
(
    id             INT AUTO_INCREMENT
        PRIMARY KEY,
    fromNumber     VARCHAR(15)  NOT NULL,
    toNumber       VARCHAR(15)  NOT NULL,
    direction      VARCHAR(20)  NOT NULL,
    updatedDate    DATETIME     NOT NULL,
    price          VARCHAR(14)  NULL,
    uri            VARCHAR(200) NOT NULL,
    accountSid     VARCHAR(34)  NOT NULL,
    status         VARCHAR(18)  NOT NULL,
    callSid        VARCHAR(34)  NOT NULL,
    sentDate       DATETIME     NOT NULL,
    createdDate    DATETIME     NOT NULL,
    priceUnit      VARCHAR(3)   NULL,
    apiVersion     VARCHAR(10)  NULL,
    parentCallSid  VARCHAR(34)  NULL,
    toFormatted    VARCHAR(20)  NULL,
    fromFormatted  VARCHAR(20)  NULL,
    phoneNumberSid VARCHAR(34)  NULL,
    answeredBy     VARCHAR(20)  NULL,
    forwardedFrom  VARCHAR(20)  NULL,
    groupSid       VARCHAR(34)  NULL,
    callerName     VARCHAR(50)  NULL,
    queueTime      VARCHAR(10)  NULL,
    trunkSid       VARCHAR(34)  NULL,
    CONSTRAINT callSid
        UNIQUE (callSid DESC)
);

CREATE TABLE sms
(
    id            INT AUTO_INCREMENT
        PRIMARY KEY,
    fromNumber    VARCHAR(15)   NOT NULL,
    toNumber      LONGTEXT      NOT NULL,
    body          VARCHAR(1600) NOT NULL,
    direction     VARCHAR(20)   NOT NULL,
    updatedDate   DATETIME      NOT NULL,
    price         VARCHAR(14)   NULL,
    uri           VARCHAR(200)  NOT NULL,
    accountSid    VARCHAR(34)   NOT NULL,
    mediaNumber   INT           NOT NULL,
    status        VARCHAR(18)   NOT NULL,
    messageSid    VARCHAR(34)   NOT NULL,
    sentDate      DATETIME      NOT NULL,
    createdDate   DATETIME      NOT NULL,
    priceUnit     VARCHAR(3)    NULL,
    apiVersion    VARCHAR(10)   NULL,
    segmentNumber INT           NOT NULL,
    errorMessage  TEXT          NULL,
    errorCode     INT           NULL,
    CONSTRAINT messageSid
        UNIQUE (messageSid DESC)
);

CREATE TABLE phoneNumberLookups
(
    id                  INT AUTO_INCREMENT
        PRIMARY KEY,
    callingCountryCode  VARCHAR(5)   NULL,
    countryCode         VARCHAR(5)   NULL,
    phoneNumber         VARCHAR(20)  NULL,
    nationalFormat      VARCHAR(20)  NULL,
    valid               TINYINT(1)   NULL,
    callerName          VARCHAR(50)  NULL,
    callerType          VARCHAR(30)  NULL,
    url                 VARCHAR(255) NULL
);