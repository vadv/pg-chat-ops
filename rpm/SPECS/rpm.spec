%define version unknown
%define bin_name pg-chat-ops
%define debug_package %{nil}

Name:           %{bin_name}
Version:        %{version}
Release:        1%{?dist}
Summary:        PlayList generator
License:        BSD
URL:            http://git.itv.restr.im/infra/%{bin_name}
Source:         %{bin_name}-%{version}.tar.gz

%define restream_dir /opt/restream/
%define restream_bin_dir %{restream_dir}/%{bin_name}/bin
%define restream_doc_dir %{restream_dir}/%{bin_name}/share/doc

%description
This package provides pg chat bot

%prep
%setup

%build
make

%install
%{__mkdir} -p %{buildroot}%{restream_bin_dir}
%{__mkdir} -p %{buildroot}%{restream_doc_dir}
%{__install} -m 0755 -p bin/%{bin_name} %{buildroot}%{restream_bin_dir}
# plugins
%{__mkdir} -p %{buildroot}%{_sysconfdir}/%{bin_name}/plugins/common
cp -v examples/plugins/*.lua %{buildroot}%{_sysconfdir}/%{bin_name}/plugins/
cp -v examples/plugins/common/*.lua %{buildroot}%{_sysconfdir}/%{bin_name}/plugins/common/
%{__install} -m 0644 examples/init.lua %{buildroot}%{_sysconfdir}/%{bin_name}/init.lua

%files
%defattr(-,root,root,-)
%{restream_bin_dir}/%{bin_name}
%{_sysconfdir}/%{bin_name}
