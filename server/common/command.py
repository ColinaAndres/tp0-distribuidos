
class Command:
    """Base class for commands that can be executed on the server."""
    def execute(self, _server, _agency):
        pass

class BetsProcessingCommand(Command):
    """Command to process bets received from a client."""
    def __init__(self, bets):
        super().__init__()
        self._bets = bets
        self._agency_id = bets[0].agency if bets else None

    def execute(self, server, agency):
        agency.agency_id = self._agency_id
        server.process_bets(self._bets, agency)

class FinalizationCommand(Command):
    """Command to indicate the finalization of the bet sending."""
    def execute(self, server, _agency):
        server.finalize_reception_of_bets()

class WinnersRequestCommand(Command):
    """Command to request the winners of the bet."""
    def execute(self, server, agency):
        server.send_winners(agency)