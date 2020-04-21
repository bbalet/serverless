CREATE DATABASE IF NOT EXISTS openfass;

USE openfass;

--
-- Database: `decouvricwp`
--

CREATE TABLE `counter` (
  `value` int(11) NOT NULL
);

INSERT INTO `counter` (`value`) VALUES
(1);
COMMIT;
