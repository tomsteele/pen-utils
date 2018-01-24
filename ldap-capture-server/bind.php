<?php
   // Test with some cool php!
   $username = "tom@tom.com";
   $password = "f00ba4";
   $ldapconn = ldap_connect("127.0.0.1");
   if (ldap_bind($ldapconn, $username, $password)) {
      print "bind successful\n";
   } else {
     print "bind failed\n";
   }
?>
