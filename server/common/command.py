
class Command:
    """Base class for commands that can be executed on the server."""
    def execute(self, _lottery_central, _agency):
        pass

class BetsProcessingCommand(Command):
    """Command to process bets received from a client."""
    def __init__(self, bets):
        super().__init__()
        self._bets = bets
        self._agency_id = bets[0].agency if bets else None

    def execute(self, lottery_central, agency):
        agency.agency_id = self._agency_id
        lottery_central.process_bets(self._bets)
        agency.send_confirmation()

class FinalizationCommand(Command):
    """Command to indicate the finalization of the bet sending."""
    def execute(self, lottery_central, _agency):
        lottery_central.register_agency_done()

class WinnersRequestCommand(Command):
    """Command to request the winners of the bet."""
    def execute(self, lottery_central, agency) -> bool:
        winners = lottery_central.get_winners_for_agency(agency.agency_id)
        agency.send_winners(winners)
        