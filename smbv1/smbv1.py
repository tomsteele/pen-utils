from impacket.smbconnection import *
import sys

host = sys.argv[1]

smb = SMBConnection(host, host, sess_port=445, preferredDialect=SMB_DIALECT)
if smb.getDialect() == SMB_DIALECT:
    print "SMBv1 suppported"
else:
    print "SMBv1 not supported"
