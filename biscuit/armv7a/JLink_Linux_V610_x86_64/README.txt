
Requirements
============
libusb is no longer required.

Running JLinkExe (or J-Link via USB in general) with standard user rights
=========================================================================
In order to run JLinkExe with standard user rights you have to do the following:

- Copy the file "99-jlink.rules" provided with this J-Link software package 
  in the /etc/udev/rules.d/ directory using this command:
  
  sudo cp 99-jlink.rules /etc/udev/rules.d/
  
  Note: For older systems it might be necessary to replace the "ATTRS" calls in the 99-jlink.rules by "SYSFS" calls

- Restart your system
