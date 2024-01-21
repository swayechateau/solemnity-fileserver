-- Migration: Initial Schema
-- Created by: Solemnity Software
-- Created on: 2023-12-26 04:56:00
-- Last Modified: 2023-12-26 04:56:00
-- Version: 0.9.0
-- Description: This migration deletes tables created by the initial schema for the Solemnity File Server.

-- Revert the addition of the Global Organizations table

-- Revert the addition of the Files table
DROP TABLE IF EXISTS Files;
-- Revert the addition of the File Access table
DROP TABLE IF EXISTS FileAccess;
-- Revert the addition of the File Recovery table
DROP TABLE IF EXISTS FileRecovery;
-- Revert the addition of the User Files table
DROP TABLE IF EXISTS UserFiles;
-- Revert the addition of the Organization Files table
DROP TABLE IF EXISTS OrganizationFiles;

