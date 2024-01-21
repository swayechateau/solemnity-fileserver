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

-- Create the Users table
CREATE TABLE IF NOT EXISTS Files (
    Id UUID PRIMARY KEY,
    FileHash TEXT NOT NULL UNIQUE, -- This is the file hash 
    FileName TEXT NOT NULL, -- This is the encrypted file name
    FilePath TEXT NOT NULL,
    FileFullPath TEXT NOT NULL,
    FileSize BIGINT NOT NULL,
    FileType TEXT NOT NULL,
    FileExtension TEXT NOT NULL,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER UpdateFilesModtime BEFORE UPDATE ON Files FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the File Access table
CREATE TABLE IF NOT EXISTS FileAccess (
    Id UUID PRIMARY KEY,
    Organization TEXT NOT NULL,
    AccessOwner TEXT,
    AccessType TEXT NOT NULL,
    Public BOOLEAN NOT NULL,
    Uri TEXT NOT NULL, -- This is the file uri
    ShareCode TEXT NOT NULL, -- This is the file share code
    AccessCode TEXT NOT NULL, -- This is the file access code
    FileId UUID,
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE,
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER UpdateFileAccessModtime BEFORE UPDATE ON FileAccess FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the File Recovery table
CREATE TABLE IF NOT EXISTS FileRecovery(
    EmailHash TEXT PRIMARY KEY, -- This is the email address hash
    EmailAddress TEXT NOT NULL, -- This is the encrypted email address
    IpAddress TEXT NOT NULL, -- This is the encrypted ip address
    Domain TEXT NOT NULL, -- This is the encrypted domain
    Code TEXT NOT NULL, -- This is the encrypted recovery code
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER UpdateFileRecoveryModtime BEFORE UPDATE ON FileRecovery FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the User Files table
CREATE TABLE IF NOT EXISTS UserFiles (
    UserId UUID,
    FileId UUID,
    FileName TEXT NULL, -- This is the encrypted file name    
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (UserId, FileId),
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE
);
CREATE TRIGGER UpdateUserFilesModtime BEFORE UPDATE ON UserFiles FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();

-- Create the Organization Files table
CREATE TABLE IF NOT EXISTS OrganizationFiles (
    OrganizationId UUID NULL,
    FileId UUID,
    FileName TEXT NULL, -- This is the encrypted file name
    CreatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (OrganizationId, FileId),
    FOREIGN KEY (FileId) REFERENCES Files(Id) ON DELETE CASCADE
);
CREATE TRIGGER UpdateOrganizationFilesModtime BEFORE UPDATE ON OrganizationFiles FOR EACH ROW EXECUTE FUNCTION UpdateUpdatedAt();
