-- Migration: Initial Schema
-- Created by: Solemnity Software
-- Created on: 2024-01-18 04:52:00
-- Last Modified: 2024-01-23 05:56:00
-- Version: 0.11.0
-- Description: This migration creates the initial schema for the Solemnity File Server.

CREATE OR REPLACE FUNCTION UpdateUpdatedAt()
RETURNS TRIGGER AS $$
BEGIN
  NEW.UpdatedAt = CURRENT_TIMESTAMP;
  RETURN NEW; 
END;
$$ language 'plpgsql';

-- organisations are the folders that files are stored in
-- Create the Logs table
CREATE TABLE IF NOT EXISTS Logs (
    Id UUID PRIMARY KEY, -- GUID
    Organization TEXT NOT NULL DEFAULT 'global', -- This is the organization name
    LogType TEXT NOT NULL, -- This is the log type
    LogMessage TEXT NOT NULL, -- This is the log message
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the Files table
CREATE TABLE IF NOT EXISTS Files (
    Id UUID PRIMARY KEY, -- GUID
    Organization TEXT NOT NULL DEFAULT 'global', -- This is the organization name
    IsEncrypted BOOLEAN NOT NULL DEFAULT FALSE, -- Is the file encrypted
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

-- Create the File Encryption table
CREATE TABLE IF NOT EXISTS FileEncryption (
    FileId UUID UNIQUE, -- File GUID
    EncryptionKey TEXT NOT NULL, -- This is the encryption key identifier
    EncryptionIterations INTEGER NOT NULL, -- Times the file was encrypted
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the File Access table
CREATE TABLE IF NOT EXISTS FileAccess (
    Id UUID PRIMARY KEY,
    FileId UUID, -- File GUID
    Organization TEXT NOT NULL DEFAULT 'global', -- This is the organization name
    AccessOwner TEXT, -- This is the Identifier for thr owner of the file
    IsPublic BOOLEAN NOT NULL,
    Uri TEXT NOT NULL UNIQUE, -- This is the file uri
    ShareCode TEXT NOT NULL, -- This is the file share code
    AccessCode TEXT NOT NULL, -- This is the file access code
    UNIQUE (Organization, AccessOwner),
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER UpdateFileAccessModtime BEFORE UPDATE ON FileAccess FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the Files Schedules for deletion table
CREATE TABLE IF NOT EXISTS FileDeletionSchedule (
    FileId UUID UNIQUE, -- File GUID
    DeleteAt TIMESTAMP WITH TIME ZONE NOT NULL, -- This is the date the file will be deleted
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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
