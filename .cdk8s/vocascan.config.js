/**
 * For more help with the file see https://docs.vocascan.com/#/vocascan-server/configuration
 * This config file lists all available options
 */

module.exports = {
  debug: false,

  server: {
    port: 8000,
    jwt_secret: '5080501ad2d4ecbb0501a0f58f66ee5d27fc1fff434bd058effc03cb83a0f40c3c1fceb51a735ec6aa6be842b72b3b400906ec77384f190a57ab0d64938fe34e',
    salt_rounds: 10,
    cors: ["*"],
  },

  database: {
    dialect: 'postgres',
    host: 'vocascan-db-service',
    port: '5432',
    username: 'postgres',
    password: 'postgres',
    database: 'vocascan',
  },

  log: {
    console: {
      level: 'info',
      colorize: true,
      enable_sql_log: true,
      enable_router_log: true,
      stderr_levels: ['error'],
    },
  },
};
