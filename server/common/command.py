
class Command:
    """Base class for commands that can be executed on the server."""
    def execute(self, server):
        pass

class BetsProcessingCommand(Command):
    """Command to process bets received from a client."""
    def __init__(self, bets):
        super().__init__()
        self._bets = bets

    def execute(self, server):
        server.process_bets(self._bets)

class FinalizationCommand(Command):
    """Command to indicate the finalization of the bet sending."""
    def execute(self, server):
        server.finalize_reception_of_bets()