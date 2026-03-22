
class Command:
    """Base class for commands that can be executed on the server."""
    def execute(self, _lottery_central, _agency) -> bool:
        pass

class BetsProcessingCommand(Command):
    """Command to process bets received from a client."""
    def __init__(self, bets):
        super().__init__()
        self._bets = bets
        self._agency_id = bets[0].agency if bets else None

    def execute(self, lottery_central, agency) -> bool:
        agency.agency_id = self._agency_id
        lottery_central.process_bets(self._bets, agency)
        return False

class FinalizationCommand(Command):
    """Command to indicate the finalization of the bet sending. returns 
    True to indicate that the session should be closed."""
    def execute(self, lottery_central, _agency) -> bool:
        lottery_central.register_agency_done()
        return True

class WinnersRequestCommand(Command):
    """Command to request the winners of the bet."""
    def execute(self, lottery_central, agency) -> bool:
        lottery_central.send_winners(agency)
        return False