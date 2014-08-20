i3event
=======

Small daemon which can listen to and handle events from I3.

Running
=======

    i3event [-config <path>] [-debug]

      -config="~/.i3event":    Path to config file.
      -debug=false:            Activate debug logging.

Configuration
=============

i3event requires a config-file which configures the rules that
incoming events are matched against. The path to the configuration
file is specified with the -config parameter. By default it will look
in "~/i3event".

The format of the configuration file looks consists of one rule per
line like this:

    bindevent <eventtype> <change> <action> [args_to_action..]

Lines starting with # is treated as comments.

Example:

    # Execute my_event_handler.sh for every window event with change:focus.
    bindevent window focus exec my_event_handler.sh