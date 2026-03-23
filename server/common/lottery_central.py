from common.utils import store_bets, load_bets, has_won
import logging
import threading

class Lottery_central:
    """
    Loterry central class taht emulates the loterry sistem and act as a coordinator
    """
    def __init__(self, total_agencies):
        self._total_agencies = total_agencies
        self._winners = []
        self._lock = threading.Lock()
        self._barrier = threading.Barrier(total_agencies, action=self._run_lottery)

    def process_bets(self, bets):
        """
        Process the bets received from an agency, storing them and logging the action,
        It ensures that the storage of bets is thread safe
        """
        with self._lock:
            store_bets(bets)
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def register_agency_done(self):
        """
        Register that an agency has finished sending bets, 
        if all agencies are done, runs the lottery
        """
        logging.info(f"action: finalizacion_recepcion_apuestas | result: success")
        self._barrier.wait()

    def get_winners_for_agency(self, agency_id) -> list:
        """
        Get the winners for a specific agency id
        As it is a lecture only operation, it does not require locking
        as the winners list is only modified in the _run_lottery method, 
        which is called after all agencies have registered as done
        """
        return list(filter(lambda bet: bet.agency == agency_id, self._winners))

    def _run_lottery(self):
        """
        Run the lottery, storing the winners and logging the action
        """
        logging.info(f"action: sorteo | result: success")
        self._winners = list(filter(has_won, load_bets()))
        