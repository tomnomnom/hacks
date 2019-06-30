<?php
/**
 * The base configuration for WordPress
 *
 * The wp-config.php creation script uses this file during the
 * installation. You don't have to use the web site, you can
 * copy this file to "wp-config.php" and fill in the values.
 *
 * This file contains the following configurations:
 *
 * * MySQL settings
 * * Secret keys
 * * Database table prefix
 * * ABSPATH
 *
 * @link https://codex.wordpress.org/Editing_wp-config.php
 *
 * @package WordPress
 */

// ** MySQL settings - You can get this info from your web host ** //
/** The name of the database for WordPress */
define( 'DB_NAME', 'wordpress' );

/** MySQL database username */
define( 'DB_USER', 'root' );

/** MySQL database password */
define( 'DB_PASSWORD', '' );

/** MySQL hostname */
define( 'DB_HOST', 'localhost' );

/** Database Charset to use in creating database tables. */
define( 'DB_CHARSET', 'utf8mb4' );

/** The Database Collate type. Don't change this if in doubt. */
define( 'DB_COLLATE', '' );

/**#@+
 * Authentication Unique Keys and Salts.
 *
 * Change these to different unique phrases!
 * You can generate these using the {@link https://api.wordpress.org/secret-key/1.1/salt/ WordPress.org secret-key service}
 * You can change these at any point in time to invalidate all existing cookies. This will force all users to have to log in again.
 *
 * @since 2.6.0
 */
define( 'AUTH_KEY',         '>#xt43EFC!s$N>Q$j}dTpw,cPgIR&BtkNwr!;C+6USvQX403J<%P4{GY*L?niN0I' );
define( 'SECURE_AUTH_KEY',  'Fms?_6sV[/)lotl$gd26k-F[wa9i#XW-Lr^;b$]h#r74@&aT*8j(ThLt1iA*XX:W' );
define( 'LOGGED_IN_KEY',    'C#n2 +T9VWs.HBqh0Xj=)zn4SdygE[Q)`(xK64oP%*>H,d}}z *p-wDX(k64d%D$' );
define( 'NONCE_KEY',        '.E9PqG<ciK_g{Md8+ <JYPw;X<yx!`e<=XhJR!>=^Y}Pf4TC67ip,VcK?L%a:g%{' );
define( 'AUTH_SALT',        ';=SYaJSSJR<g Uc(cwNsy!Rr5TU5h-:40`p:mr{hP$T@mq(yq@c0Bai_&<ml+[]9' );
define( 'SECURE_AUTH_SALT', '.Wf4rHeY)$nh.#o+OwI C7WVX6#`?>c;wai_s78um/alrf7F.pz$l!_(NBiS=I9Z' );
define( 'LOGGED_IN_SALT',   'bY74>41pZopfHmvkT~VdC1|.|KSYgQ|mh;n(j <@mL3Sx,(=Nl|5nv5]ZWu)K^(U' );
define( 'NONCE_SALT',       'N6Bs8YQ6?H1QAY?jWBW]zn|.Rd4%Z:n=UaT#*ech0sR4sU<176IT(Nb1I5E/>d@3' );

/**#@-*/

/**
 * WordPress Database Table prefix.
 *
 * You can have multiple installations in one database if you give each
 * a unique prefix. Only numbers, letters, and underscores please!
 */
$table_prefix = 'wp_';

/**
 * For developers: WordPress debugging mode.
 *
 * Change this to true to enable the display of notices during development.
 * It is strongly recommended that plugin and theme developers use WP_DEBUG
 * in their development environments.
 *
 * For information on other constants that can be used for debugging,
 * visit the Codex.
 *
 * @link https://codex.wordpress.org/Debugging_in_WordPress
 */
define( 'WP_DEBUG', false );

// direct plugin installs
define('FS_METHOD', 'direct');

/* That's all, stop editing! Happy publishing. */

/** Absolute path to the WordPress directory. */
if ( ! defined( 'ABSPATH' ) ) {
	define( 'ABSPATH', dirname( __FILE__ ) . '/' );
}

/** Sets up WordPress vars and included files. */
require_once( ABSPATH . 'wp-settings.php' );

