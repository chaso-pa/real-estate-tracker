-- CreateTable
CREATE TABLE `estates` (
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `url` VARCHAR(191) NOT NULL,
    `address` VARCHAR(191) NULL,
    `estate_type` ENUM('used_house', 'new_house', 'land', 'used_apartment', 'new_apartment') NULL,
    `value` INTEGER NULL,
    `railway` VARCHAR(191) NULL,
    `land_area` DOUBLE NULL,
    `building_area` DOUBLE NULL,
    `floor_plan` VARCHAR(191) NULL,
    `year_of_construction` INTEGER NULL,
    `first_appeared` DATETIME(3) NULL,
    `last_appeared` DATETIME(3) NULL,
    `created_date` DATETIME(3) NULL,
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) NOT NULL,

    UNIQUE INDEX `estates_url_key`(`url`),
    PRIMARY KEY (`id`)
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
