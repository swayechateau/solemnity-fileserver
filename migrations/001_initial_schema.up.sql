-- Migration: Initial Schema
-- Created by: SolemnitySSO
-- Created on: 2023-12-26 04:52:00
-- Last Modified: 2023-01-02 05:56:00
-- Version: 0.11.0
-- Description: This migration creates the initial schema for the SolemnitySSO OAuth Server.

CREATE OR REPLACE FUNCTION UpdateUpdatedAt()
RETURNS TRIGGER AS $$
BEGIN
  NEW.UpdatedAt = CURRENT_TIMESTAMP;
  RETURN NEW; 
END;
$$ language 'plpgsql';

-- Create the Files table
CREATE TABLE IF NOT EXISTS Files (
    Id UUID PRIMARY KEY, -- GUID
    Organization TEXT NOT NULL DEFAULT 'global', -- This is the organization name
    FileHash TEXT NOT NULL, -- This is the file hash 
    FileName TEXT NOT NULL, -- This is the encrypted file name
    FilePath TEXT NOT NULL,
    FileFullPath TEXT NOT NULL,
    FileSize BIGINT NOT NULL,
    FileType TEXT NOT NULL,
    FileExtension TEXT NOT NULL,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (Organization, FileHash)
);
CREATE TRIGGER UpdateFilesModtime BEFORE UPDATE ON Files FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- -- Create file File Metadata table
-- CREATE TABLE IF NOT EXISTS FileMetas (
--     FileId UUID UNIQUE, -- File GUID
--     FileName TEXT NOT NULL, -- This is the encrypted file name
--     FilePath TEXT NOT NULL,
--     FileFullPath TEXT NOT NULL,
--     FileSize BIGINT NOT NULL,
--     FileType TEXT NOT NULL,
--     FileExtension TEXT NOT NULL,
--     CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE
-- );
-- CREATE TRIGGER UpdateFileMetasModtime BEFORE UPDATE ON FileMetas FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the File Access table
CREATE TABLE IF NOT EXISTS FileAccess (
    Id UUID PRIMARY KEY,
    FileId UUID, -- File GUID
    Organization TEXT NOT NULL DEFAULT 'global', -- This is the organization name
    AccessOwner TEXT, -- This is the Identifier for thr owner of the file
    Public BOOLEAN NOT NULL,
    Uri TEXT NOT NULL UNIQUE, -- This is the file uri
    ShareCode TEXT NOT NULL, -- This is the file share code
    AccessCode TEXT NOT NULL, -- This is the file access code
    UNIQUE (Organization, AccessOwner),
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER UpdateFileAccessModtime BEFORE UPDATE ON FileAccess FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- -- Create the File Recovery table
-- CREATE TABLE IF NOT EXISTS FileRecovery(
--     EmailHash TEXT PRIMARY KEY, -- This is the email address hash
--     EmailAddress TEXT NOT NULL, -- This is the encrypted email address
--     IpAddress TEXT NOT NULL, -- This is the encrypted ip address
--     Domain TEXT NOT NULL, -- This is the encrypted domain
--     Code TEXT NOT NULL, -- This is the encrypted recovery code
--     CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
-- );
-- CREATE TRIGGER UpdateFileRecoveryModtime BEFORE UPDATE ON FileRecovery FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();
