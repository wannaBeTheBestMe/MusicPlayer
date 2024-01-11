DROP TABLE IF EXISTS Tracks;
DROP TABLE IF EXISTS album;

CREATE TABLE Tracks
(
	id                 INT AUTO_INCREMENT NOT NULL,
	Format             VARCHAR(50)        NOT NULL,
	FileType           VARCHAR(50)        NOT NULL,
	Title              VARCHAR(255)       NOT NULL,
	Album              VARCHAR(255),
	Artist             VARCHAR(255),
	AlbumArtist        VARCHAR(255),
	Composer           VARCHAR(255),
	Year               INT,
	Genre              VARCHAR(100),
	TrackNum           INT,
	TrackTotal         INT,
	DiscNum            INT,
	DiscTotal          INT,
	PictureExt         VARCHAR(50),
	PictureMIMEType    VARCHAR(255),
	PictureType        VARCHAR(255),
	PictureDescription TEXT,
	PictureData        MEDIUMBLOB,
	Lyrics             TEXT,
	Comment            TEXT,
	PRIMARY KEY (`id`)
);

CREATE TABLE album
(
	id     INT AUTO_INCREMENT NOT NULL,
	title  VARCHAR(128)       NOT NULL,
	artist VARCHAR(255)       NOT NULL,
	price  DECIMAL(5, 2)      NOT NULL,
	PRIMARY KEY (`id`)
);

INSERT INTO album
(title, artist, price)
VALUES ('Blue Train', 'John Coltrane', 56.99),
			 ('Giant Steps', 'John Coltrane', 63.99),
			 ('Jeru', 'Gerry Mulligan', 17.99),
			 ('Sarah Vaughan', 'Sarah Vaughan', 34.98);