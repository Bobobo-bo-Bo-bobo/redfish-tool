Name:           redfish-tool
Version:        1.2.1
Release:        1%{?dist}
Summary:        Tool to perform various operations on management boards using the Redfish REST API

Group:          System Environment/Daemons
License:        GPL
URL:            https://git.ypbind.de/cgit/redfish-tool
Source0:        https://git.ypbind.de/cgit/redfish-tool/snapshot/redfish-tool-1.2.1.tar.gz
BuildRequires:  golang

# don't build debuginfo package
%define debug_package %{nil}

# Don't barf on missing build_id
%global _missing_build_ids_terminate_build 0

%description
Tool to perform various operations on servers and management boards using the Redfish REST API.

%prep
%setup -q


%build
make depend build strip


%install
make install DESTDIR=%{buildroot}

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/redfish-tool

%changelog
* Sat Jul 04 2020 Andreas Maus <andreas.maus@atos.net> - 1.2.1
- support HP(E) privilege maps for users
- allow for simplified options for CSR generation (-o/-ou/...)
- add support for HP(E) iLO licenses
- for files to read, support "-" for stdin
- output can be formatted as text or JSON
- fix password handling of read from terminal and ends with
  whitespace (Issue#6)
- show version and use semantic versioning

